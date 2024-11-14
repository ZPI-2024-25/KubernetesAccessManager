// import styles from './LoginForm.module.css';
// import {Button, Form, Input} from "antd";
// import {GoogleOutlined, LockOutlined, UserOutlined} from "@ant-design/icons";
// import React from "react";
// import { Typography } from 'antd';
// const { Title } = Typography;
//
//
// const LoginForm: React.FC = () => {
//     const onFinish = (values: never) => {
//         console.log('Success:', values);
//     };
//
//     const onFinishFailed = (errorInfo: never) => {
//         console.log('Failed:', errorInfo);
//     };
//     return (
//         <div className={styles.loginContainer}>
//             <div className={styles.loginBox}>
//                 <Title level={3} className={styles.loginTitle}>Log In</Title>
//                 <Form
//                     name="login"
//                     initialValues={{ remember: true }}
//                     onFinish={onFinish}
//                     onFinishFailed={onFinishFailed}
//                 >
//                     <Form.Item
//                         name="username"
//                         rules={[{ required: true, message: 'Please input your username!' }]}
//                     >
//                         <Input className={styles.input} prefix={<UserOutlined />} placeholder="Username" />
//                     </Form.Item>
//
//                     <Form.Item
//                         name="password"
//                         rules={[{ required: true, message: 'Please input your password!' }]}
//                     >
//                         <Input.Password className={styles.input} prefix={<LockOutlined />} placeholder="Password" />
//                     </Form.Item>
//
//                     <Form.Item>
//                         <a className={styles.loginForgot} href="/forgot-password">
//                             Forgot password?
//                         </a>
//                     </Form.Item>
//                     <Form.Item>
//                         <Button type="primary" htmlType="submit" className={styles.loginButton}>
//                             Log In
//                         </Button>
//                     </Form.Item>
//
//                     <Form.Item>
//                         <Button type="default" className={styles.googleLogin}>
//                             <GoogleOutlined /> Log In with Google
//                         </Button>
//                     </Form.Item>
//                     <Form.Item>
//                         <Button type="default" className={styles.ssoLogin}>
//                             Log In with SSO
//                         </Button>
//                     </Form.Item>
//
//                     <Form.Item>
//                         <span>Don't have an account?</span> <a href="/register">Register</a>
//                     </Form.Item>
//                 </Form>
//             </div>
//         </div>
//     );
// };
// export default LoginForm;