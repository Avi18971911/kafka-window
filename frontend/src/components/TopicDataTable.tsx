import React from "react";
import TopicDataTableCell from "./TopicDataTableCell.tsx";
import styles from "../styles/DataTable.module.css";
import {TopicDetails} from "../model/TopicDetails.ts";

type TopicDataTableProps = {
    topics: TopicDetails[]
}

const TopicDataTable: React.FC<TopicDataTableProps> = ({ topics: topicDetails }) => {
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
                        <TopicDataTableCell
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

export default TopicDataTable;