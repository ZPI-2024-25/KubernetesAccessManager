import React, { useState, useEffect } from 'react';
import { Breadcrumb, Layout, Menu } from 'antd';
import styles from './Menu.module.css';
import { items } from '../../consts/MenuItem';
import { MenuItem } from '../../types';
import { Outlet } from 'react-router-dom';

const { Header, Content, Footer, Sider } = Layout;

const LeftMenu: React.FC = () => {
    const [collapsed, setCollapsed] = useState<boolean>(false);
    const [currentPage, setCurrentPage] = useState<string>('');
    const [selectedKey, setSelectedKey] = useState<string>('1');
    const [asideWidth, setAsideWidth] = useState<number>(250);
    const username = 'k8_userjjjjjjjjjjjjjjjiiiiiiiii';

    const findItemByKey = (key: string, items: MenuItem[]): MenuItem | { sectionLabel: string, childLabel: string } | null => {
        for (const item of items) {
            if (item.key === key) {
                return item;
            }
            // If the item has children, search recursively in children
            if (item.children) {
                const found = findItemByKey(key, item.children);
                if (found) {
                    return {
                        sectionLabel: item.label,
                        childLabel: (found as MenuItem).label,
                    };
                }
            }
        }
        return null;
    };

    const setCurrentPageFromItem = (item: MenuItem | { sectionLabel: string, childLabel: string }) => {
        if ('childLabel' in item) {
            setCurrentPage(`${item.sectionLabel}/${item.childLabel}`);
        } else {
            setCurrentPage(item.label);
        }
    };

    const handleMenuClick = (e: { key: string }) => {
        setSelectedKey(e.key);
        const selectedItem = findItemByKey(e.key, items);

        if (selectedItem) {
            setCurrentPageFromItem(selectedItem);
        } else {
            setCurrentPage('');
        }
    };

    // Effect to set the default page title based on the initially selected key
    useEffect(() => {
        const defaultItem = findItemByKey(selectedKey, items);
        if (defaultItem) {
            setCurrentPageFromItem(defaultItem);
        }
    }, []);

    return (
        <Layout className={styles.menuLayout}>
            <Sider
                className={styles.menuSider}
                collapsible
                collapsed={collapsed}
                onCollapse={(value) => {
                    setCollapsed(value);
                    setAsideWidth(asideWidth === 80 ? 250 : 80);
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
                    selectedKeys={[selectedKey]}
                    mode="inline"
                    items={items}
                    onClick={handleMenuClick}
                />
            </Sider>
            <Layout className={styles.contentLayout} style={{marginLeft: `${asideWidth}px`}}>
                <Header className={styles.header}><p style={{paddingLeft: `${asideWidth}px`}}>{currentPage || 'name of page'}</p></Header>
                <Content className={styles.content}>
                    <Breadcrumb className={styles.breadcrumb}>
                        <Breadcrumb.Item>lists</Breadcrumb.Item>
                    </Breadcrumb>
                    <div className={styles.innerContent}>
                        <Outlet />
                    </div>
                </Content>
                <Footer className={styles.footer}>
                    ZPI Kubernetes Access Manager ©{new Date().getFullYear()} Created by SDVM
                </Footer>
            </Layout>
        </Layout>
    );
};

export default LeftMenu;
