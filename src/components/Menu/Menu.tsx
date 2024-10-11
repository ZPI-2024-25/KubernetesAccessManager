// LeftMenu.tsx
import React, { useState } from 'react';
import {
    DesktopOutlined,
    FileOutlined,
    PieChartOutlined,
    TeamOutlined,
    UserOutlined,
} from '@ant-design/icons';
import type { MenuProps } from 'antd';
import { Breadcrumb, Layout, Menu, theme } from 'antd';
import styles from './Menu.module.css'; // Импортируем стили

const { Header, Content, Footer, Sider } = Layout;

type MenuItem = Required<MenuProps>['items'][number];

function getItem(
    label: React.ReactNode,
    key: React.Key,
    icon?: React.ReactNode,
    children?: MenuItem[],
): MenuItem {
    return {
        key,
        icon,
        children,
        label,
    } as MenuItem;
}

const items: MenuItem[] = [
    getItem('Nodes', '1', <PieChartOutlined />),
    getItem('Workloads', 'sub1', <UserOutlined />, [
        getItem('Overview', '3'),
        getItem('Pods', '4'),
        getItem('Deployments', '5'),
        getItem('StatefulSets', '6'),
        getItem('ReplicaSets', '7'),
    ]),
    getItem('Config', 'sub2', <TeamOutlined />, [
        getItem('Config 1', '9'),
        getItem('Config 2', '10'),
    ]),
    getItem('Network', 'sub3', <TeamOutlined />, [
        getItem('Network 1', '12'),
        getItem('Network 2', '13'),
    ]),
    getItem('Storage', 'sub4', <TeamOutlined />, [
        getItem('Storage 1', '15'),
        getItem('Storage 2', '16'),
    ]),
    getItem('Namespaces', '17', <FileOutlined />),
    getItem('Events', '18', <FileOutlined />),
    getItem('Access Control', 'sub5', <TeamOutlined />, [
        getItem('Access Control 1', '20'),
        getItem('Access Control 2', '21'),
    ]),
];

const LeftMenu: React.FC = () => {
    const [collapsed, setCollapsed] = useState(false);
    const {
        token: { colorBgContainer, borderRadiusLG },
    } = theme.useToken();

    return (
        <Layout style={{ minHeight: '100vh' }}>
            <Sider collapsible collapsed={collapsed} onCollapse={(value) => setCollapsed(value)}>
                <div className={styles.logo} />
                <Menu theme="dark" defaultSelectedKeys={['1']} mode="inline" items={items} />
            </Sider>
            <Layout>
                <Header className={styles.header} />
                <Content className={styles.content}>
                    <Breadcrumb className={styles.breadcrumb}>
                        <Breadcrumb.Item>User</Breadcrumb.Item>
                        <Breadcrumb.Item>Bill</Breadcrumb.Item>
                    </Breadcrumb>
                    <div className={styles.innerContent}>
                        Bill is a cat.
                    </div>
                </Content>
                <Footer className={styles.footer}>
                    Ant Design ©{new Date().getFullYear()} Created by Ant UED
                </Footer>
            </Layout>
        </Layout>
    );
};

export default LeftMenu;
