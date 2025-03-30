import {NavLink} from "react-router-dom";
import routes from "../config/routes.ts";
import styles from "../styles/NavLink.module.css";

function Toolbar() {
    return (
        <nav style={{
            background: 'gray', position: 'sticky', top: 0, left: 0, width: '100%', height: '40px', zIndex: 100,
        }}>
            <ul style={{
                display: 'flex', listStyle: 'none', padding: 0, marginLeft: 20, marginTop: 0,
                marginBottom: 0, height: '100%'
            }}>
                {Object.entries(routes).map(([key, route]) => (
                    <li key={key} style={{ height: '100%', width: '80px'}}>
                        <NavLink
                            to={route.path}
                            className={({ isActive }) =>
                                isActive ? styles.active : styles.inactive
                            }
                        >
                            {route.label}
                        </NavLink>
                    </li>
                ))}
            </ul>
        </nav>
    );
}

export default Toolbar