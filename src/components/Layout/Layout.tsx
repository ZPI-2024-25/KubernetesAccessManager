import {ReactNode} from "react";
import styles from './Layout.module.css'
import {Outlet} from "react-router-dom";

export const Layout = () => {
    return <div className={styles.pageWrapper}>
        <div className={styles.container}>
        <Outlet/>
        </div>
    </div>
}