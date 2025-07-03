import React from "react";

export type partitionNumberOption = 'All' | number

type TopicDetailsOptionBarProps = {
    partitions: partitionNumberOption[],
    onPartitionChange: (partition: partitionNumberOption) => void
    onPartitionDetailsChange: (partition: partitionNumberOption, startOffset: number, endOffset: number) => void
}

const TopicDetailsOptionBar: React.FC<TopicDetailsOptionBarProps> = (
    { partitions, onPartitionChange, onPartitionDetailsChange }
) => {
    const handlePartitionChange = (event: React.ChangeEvent<HTMLSelectElement>) => {
        if (event.target.value === 'All') {
            onPartitionChange('All');
        } else {
            const partitionNumber = parseInt(event.target.value, 10);
            if (isNaN(partitionNumber)) {
                console.error("Invalid partition number selected:", event.target.value);
                return;
            }
            if (!partitions.includes(partitionNumber)) {
                console.error("Selected partition number is not in the list of partitions:", partitionNumber);
                return;
            }
            onPartitionChange(partitionNumber);
        }
    }

    return (
        <div
            style={{ display: 'flex', gap: '1rem'}}
        >
            <div
                style={{
                    display: 'flex',
                    gap: '0.5rem',
                    alignItems: 'center',
                    flexDirection: 'column',
                    fontSize: '0.875rem',
                }}
            >
                Partition
                <select defaultValue = 'All' onChange={handlePartitionChange}>
                    <option key='All' value='All'> All </option>
                    {
                        partitions.map((partition) => (
                            <option key={partition} value={partition}>
                                {partition}
                            </option>
                        ))
                    }
                </select>
            </div>

            <select>
                <option value="asc">Ascending</option>
                <option value="desc">Descending</option>
            </select>
        </div>
    )
}

export default TopicDetailsOptionBar;