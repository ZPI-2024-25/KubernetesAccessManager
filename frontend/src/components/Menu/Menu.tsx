import React, { useState } from 'react';
import { Layout, Menu } from 'antd';
import styles from './Menu.module.css';
import { items } from '../../consts/MenuItem';
import { MenuItem } from '../../types';
import {Link, Outlet} from "react-router-dom";
import { useLocation } from 'react-router-dom';

const { Header, Content, Footer, Sider } = Layout;

const LeftMenu: React.FC = () => {
    const [collapsed, setCollapsed] = useState<boolean>(false);
    const [asideWidth, setAsideWidth] = useState<number>(270);
    const username = 'k8_userjjjjjjjjjjjjjjjiiiiiiiii';
    const location = useLocation();

    const onCollapse = (value: boolean) => {
        setCollapsed(value);
        setAsideWidth(value ? 80 : 250);
    };

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
                    <Link to={`/${item.resourceLabel || ''}`}>
                        {item.label}
                    </Link>
                ),
            };
        });
    };

    const getSelectedKeys = (menuItems: MenuItem[], pathname: string): string[] => {
        for (const item of menuItems) {
            const itemPath = `/${item.resourceLabel || ''}`;
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
                <div className={styles.logo}>
          <span className={styles.logoText}>
            {collapsed ? 'U' : username.length > 10 ? `${username.slice(0, 10)}...` : username}
          </span>
                </div>
                <Menu
                    theme="dark"
                    selectedKeys={selectedKeys}
                    mode="inline"
                    items={generateMenuItems(items)
                    }
                    style={{ paddingBottom: 50 }}

                />
            </Sider>
            <Layout className={styles.contentLayout} style={{ marginLeft: asideWidth }}>
                <Header className={styles.header}>
                    <p style={{ paddingLeft: asideWidth }}>
                        {currentPageTitle || 'Page name'}
                    </p>
                </Header>
                <Content className={styles.content}>
                    <Outlet />
                </Content>
                <Footer className={styles.footer}>
                    ZPI Kubernetes Access Manager ©{new Date().getFullYear()} Created by SDVM
                </Footer>
            </Layout>
        </Layout>
    );
};

export default LeftMenu;