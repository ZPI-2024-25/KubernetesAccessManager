import React from "react";
import {Layout, Menu, Typography} from "antd";
import styles from './Navbar.module.css'

const { Header } = Layout;
const { Title } = Typography;

export const Navbar: React.FC = () => {
    return (
        <Layout className={styles.navbar}>
            <Header className={styles.navbarHeader}>
                <div style={{ flex: '1' }}>
                    <Title level={2} style={{ color: 'white', margin: 0 }}>LangCards</Title>
                </div>
                <Menu theme="dark" mode="horizontal" >
                    <Menu.Item key={1}>My account</Menu.Item>
                    <Menu.Item key={2}>Help</Menu.Item>
                    <Menu.Item key={3}>Log Out</Menu.Item>
                </Menu>
            </Header>
        </Layout>
    );
};
