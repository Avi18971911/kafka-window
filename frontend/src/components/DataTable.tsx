import React from "react";
import DataTableCell from "./DataTableCell.tsx";
import styles from "../styles/DataTable.module.css";
import {TopicDetails} from "../model/TopicDetails.ts";

type DataTableProps = {
    topics: TopicDetails[]
}

const DataTable: React.FC<DataTableProps> = ({ topics: topicDetails }) => {
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
                        <th>Number of Partitions</th>
                        <th>Replication Factor</th>
                        <th>Cleanup Policy</th>
                        <th>Is Internal Topic?</th>
                        <th>Retention (ms)</th>
                        <th>Retention (bytes)</th>
                    </tr>
                </thead>
                <tbody>
                    {topicDetails.map((topicDetails) => (
                        <DataTableCell
                            key={topicDetails.topic}
                            topicDetails={topicDetails}
                            onToggleExpand={toggleExpand}
                            expanded={expandedRow === topicDetails.topic}
                        />
                    ))}
                </tbody>
            </table>
        </div>
    );
}

export default DataTable;