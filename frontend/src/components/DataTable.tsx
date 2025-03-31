import React from "react";
import DataTableCell from "./DataTableCell.tsx";
import styles from "../styles/DataTable.module.css";

type DataTableProps = {
    topics: string[]
}

const DataTable: React.FC<DataTableProps> = ({ topics }) => {
    const [expandedRow, setExpandedRow] = React.useState<string | null>(null);

    const toggleExpand = (id: string) => {
        setExpandedRow((prev) => (prev === id ? null : id));
    };

    return (
        <div className={styles.tableContainer}>
            <table className={styles.table}>
                <thead>
                    <tr>
                        <th></th>
                        <th>Topic Name</th>
                    </tr>
                </thead>
                <tbody>
                    {topics.map((topic) => (
                        <DataTableCell
                            key={topic}
                            topic={topic}
                            onToggleExpand={toggleExpand}
                            expanded={expandedRow === topic}
                        />
                    ))}
                </tbody>
            </table>
        </div>
    );
}

export default DataTable;