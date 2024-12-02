import React, { useState } from 'react';
import { Button, Layout, Menu } from 'antd';
import styles from './Menu.module.css';
import { items } from '../../consts/MenuItem';
import { MenuItem } from '../../types';
import { Link, Outlet, useLocation } from "react-router-dom";
import { useAuth } from '../AuthProvider/AuthProvider';
import { hasAnyPermissionInAnyNamespace } from '../../functions/authorization';


const { Header, Content, Sider } = Layout;

const LeftMenu: React.FC = () => {
    const [collapsed, setCollapsed] = useState<boolean>(false);
    const [asideWidth, setAsideWidth] = useState<number>(270);
    const {user, isLoggedIn, handleLogout, permissions } = useAuth();
    const location = useLocation();

    const generateMenuItems = (menuItems: MenuItem[]): MenuItem[] => {
        return menuItems.map((item) => {
            if (item.children) {
                return {
                    ...item,
                    children: generateMenuItems(item.children),
                };
            }
            const disabled = permissions !== null && !hasAnyPermissionInAnyNamespace(permissions, item.resourcelabel)
            return {
                ...item,
                label: (
                    <Link to={`/${item.resourcelabel || ''}`} onClick={(e) => disabled && e.preventDefault()}>
                        {item.label}
                    </Link>
                ),
                disabled: disabled,
            };
        });
    };

    const getSelectedKeys = (menuItems: MenuItem[], pathname: string): string[] => {
        for (const item of menuItems) {
            const itemPath = `/${item.resourcelabel || ''}`;
            if (itemPath === pathname) {
                return [item.key];
            }
            if (item.children) {
                const childSelectedKeys = getSelectedKeys(item.children, pathname);
                if (childSelectedKeys.length > 0) {
                    return [item.key, ...childSelectedKeys];
                }
            }
        }
        return [];
    };

    const getCurrentPageTitleFromKeys = (menuItems: MenuItem[], keys: string[]): string => {
        const labels: string[] = [];
        let currentItems = menuItems;
        for (const key of keys) {
            const item = currentItems.find((item) => item.key === key);
            if (item) {
                labels.push(item.label as string);
                if (item.children) {
                    currentItems = item.children;
                } else {
                    break;
                }
            } else {
                break;
            }
        }
        return labels.join('/');
    };

    const selectedKeys = getSelectedKeys(items, location.pathname);
    const currentPageTitle = getCurrentPageTitleFromKeys(items, selectedKeys);

    const isLoginPage = location.pathname === "/login";

    const menuItems = isLoginPage ? [] : generateMenuItems(items);

    return (
        <Layout className={styles.menuLayout}>
            <Sider
                className={styles.menuSider}
                collapsible
                collapsed={collapsed}
                onCollapse={(value) => {
                    setCollapsed(value);
                    setAsideWidth(asideWidth === 80 ? 270 : 80);
                }}
                width={`${asideWidth}px`}
            >
                <div className={styles.user} style={{marginLeft: `${collapsed ? 10 : 40}px`, marginRight: `${collapsed ? 10 : 40}px`}}>
                    <span className={styles.userText}>
                        {collapsed ? (user?.preferred_username ? user.preferred_username.slice(0, 2) : 'U') : (user?.preferred_username || 'User')}
                    </span>
                </div>

                <Menu
                    className={styles.menuMenu}
                    theme="dark"
                    selectedKeys={selectedKeys}
                    mode="inline"
                    items={menuItems}
                />
            </Sider>
            <Layout
                className={styles.contentLayout}
                style={{marginLeft: !isLoggedIn ? '0' : `${asideWidth}px`}}
            >
                <Header className={styles.header} style={{paddingLeft: !isLoggedIn ? '0' : `${asideWidth}px`}}>
                    <div>
                        <p>{currentPageTitle || ' '}</p>

                        <Button type="primary" onClick={handleLogout}>
                            Log out
                        </Button>
                    </div>
                </Header>
                <Content
                    className={styles.content}
                    style={{
                        padding: '16px',
                        minHeight: 'calc(100vh - 64px - 50px)',
                        height: 'auto',
                    }}
                >
                    <Outlet/>
                </Content>
                {/*<Footer className={styles.footer}>*/}
                {/*    ZPI Kubernetes Access Manager ©{new Date().getFullYear()} Created by SDVM*/}
                {/*</Footer>*/}
            </Layout>
        </Layout>
    );
};

export default LeftMenu;
