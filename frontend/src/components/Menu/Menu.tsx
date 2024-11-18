import React, { useEffect, useState } from 'react';
import { Button, Layout, Menu } from 'antd';
import { jwtDecode } from 'jwt-decode';
import styles from './Menu.module.css';
import { items } from '../../consts/MenuItem';
import { MenuItem } from '../../types';
import { Link, Outlet, useLocation } from "react-router-dom";

const { Header, Content, Footer, Sider } = Layout;

const LeftMenu: React.FC = () => {
    const [collapsed, setCollapsed] = useState<boolean>(false);
    const [asideWidth, setAsideWidth] = useState<number>(270);
    const [isLoggedIn, setIsLoggedIn] = useState<boolean>(false);
    const [username, setUsername] = useState<string>('Użytkownik');
    const location = useLocation();

    // Funkcja dekodująca token JWT
    const decodeToken = (token: string): string | null => {
        try {
            const decoded: { preferred_username?: string } = jwtDecode(token);
            return decoded.preferred_username || null;
        } catch (error) {
            console.error("Nie udało się zdekodować tokena JWT:", error);
            return null;
        }
    };

    useEffect(() => {
        // Funkcja do aktualizacji stanu logowania
        const updateLoginState = () => {
            const token = localStorage.getItem('access_token');
            if (token) {
                setIsLoggedIn(true);
                const preferredUsername = decodeToken(token);
                if (preferredUsername) {
                    setUsername(preferredUsername);
                }
            } else {
                setIsLoggedIn(false);
                setUsername('Użytkownik');
            }
        };

        // Zaktualizuj stan przy każdej zmianie ścieżki (np. po powrocie z logowania)
        updateLoginState();
    }, [location.pathname]); // Nasłuchujemy na zmianę ścieżki

    const handleLogin = () => {
        const redirectUri = `${window.location.origin}/auth/callback`;
        window.location.href = `http://localhost:4000/realms/ZPI-realm/protocol/openid-connect/auth?client_id=ZPI-client&response_type=code&redirect_uri=${encodeURIComponent(redirectUri)}`;
    };

    const handleLogout = () => {
        const logoutUrl = `http://localhost:4000/realms/ZPI-realm/protocol/openid-connect/logout?redirect_uri=${encodeURIComponent(window.location.origin)}`;
        localStorage.removeItem('access_token');
        localStorage.removeItem('refresh_token');
        setIsLoggedIn(false);
        setUsername('Użytkownik');
        window.location.href = logoutUrl;
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
                    items={generateMenuItems(items)}
                    style={{ paddingBottom: 50 }}
                />
            </Sider>
            <Layout className={styles.contentLayout} style={{ marginLeft: asideWidth }}>
                <Header className={styles.header}>
                    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                        <p style={{ paddingLeft: asideWidth }}>
                            {currentPageTitle || 'Page name'}
                        </p>
                        {isLoggedIn ? (
                            <Button type="primary" onClick={handleLogout}>
                                Wyloguj
                            </Button>
                        ) : (
                            <Button type="primary" onClick={handleLogin}>
                                Zaloguj
                            </Button>
                        )}
                    </div>
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
