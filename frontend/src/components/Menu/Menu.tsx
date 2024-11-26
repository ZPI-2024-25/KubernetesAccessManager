import React, { useState } from 'react';
import { Button, Layout, Menu } from 'antd';
import styles from './Menu.module.css';
import { items } from '../../consts/MenuItem';
import { MenuItem } from '../../types';
import { Link, Outlet, useLocation } from "react-router-dom";
import { useAuth } from '../AuthProvider/AuthProvider';


const { Header, Content, Sider } = Layout;

const LeftMenu: React.FC = () => {
    const [collapsed, setCollapsed] = useState<boolean>(false);
    const [asideWidth, setAsideWidth] = useState<number>(270);
    const { user, isLoggedIn, handleLogin, handleLogout } = useAuth();
    const location = useLocation();

    const generateMenuItems = (menuItems: MenuItem[]): MenuItem[] => {
        return menuItems.map((item) => {
            if (item.children) {
                return {
                    ...item,
                    children: generateMenuItems(item.children),
                };
            }
            return {
                ...item,
                label: (
                    <Link to={`/${item.resourcelabel || ''}`}>
                        {item.label}
                    </Link>
                ),
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
                style={{ display: !isLoggedIn ? 'none' : 'block' }}
            >
                <div className={styles.logo}>
                    <span
                        className={styles.logoText}
                        title={user?.preferred_username || 'User'}
                        style={{
                            maxWidth: `${asideWidth - 20}px`,
                        }}
                    >
                        {user?.preferred_username || 'User'}
                    </span>
                </div>

                <Menu
                    theme="dark"
                    selectedKeys={selectedKeys}
                    mode="inline"
                    items={menuItems}
                    style={{ paddingBottom: 50 }}
                />
            </Sider>
            <Layout
                className={styles.contentLayout}
                style={{ marginLeft: !isLoggedIn ? '0' : `${asideWidth}px` }}
            >
                <Header className={styles.header}>
                    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', height: '100%' }}>
                        <p style={{ paddingLeft: isLoggedIn ? `${asideWidth}px` : '0'}}>
                            {currentPageTitle || 'Page name'}
                        </p>
                        {isLoggedIn ? (
                            <Button type="primary" onClick={handleLogout}>
                                Log out
                            </Button>
                        ) : (
                            <Button type="primary" onClick={handleLogin}>
                                Log in
                            </Button>
                        )}
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
                    <Outlet />
                </Content>
                {/*<Footer className={styles.footer}>*/}
                {/*    ZPI Kubernetes Access Manager ©{new Date().getFullYear()} Created by SDVM*/}
                {/*</Footer>*/}
            </Layout>
        </Layout>
    );
};

export default LeftMenu;
