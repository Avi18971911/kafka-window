import React from "react";

type DataTableCellProps = {
    topic: string;
    onToggleExpand: (topic: string) => void;
    expanded: boolean;
};

const DataTableCell: React.FC<DataTableCellProps> = ({ topic, onToggleExpand, expanded }) => {
    return (
        <>
            <tr>
                <td>
                    <button
                        onClick={() => onToggleExpand(topic)}
                        style={{
                            background: 'none',
                            border: 'none',
                            cursor: 'pointer',
                            color: 'white',
                            fontSize: '16px',
                        }}
                    >
                        {expanded ? '▼' : '▶'}
                    </button>
                </td>
                <td>{topic}</td>
            </tr>

            {expanded && (
                <tr>
                    <td colSpan={2}>
                        <pre>
                            {JSON.stringify(topic, null, 2)}
                        </pre>
                    </td>
                </tr>
            )}
        </>
    );
};

export default DataTableCell;