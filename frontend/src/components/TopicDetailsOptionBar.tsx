import React from "react";

export type partitionNumberOption = 'All' | number

type partitionOffsetOption = 'Latest' | 'Earliest' | 'Custom'

type TopicDetailsOptionBarProps = {
    partitions: partitionNumberOption[],
    onPartitionChange: (partition: partitionNumberOption) => void
    onPartitionDetailsChange: (startOffset: number, endOffset: number) => void
}

const partitionOffsetOptions: partitionOffsetOption[] = [
    'Latest',
    'Earliest',
    'Custom',
]

type offsetState = {
    startOffset: number,
    endOffset: number,
}

const defaultNumMessages = 50;
const defaultStartOffset = -1*defaultNumMessages;
const defaultEndOffset = -1;

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

    const handleOffsetChange = (event: React.ChangeEvent<HTMLSelectElement>) => {
        const selectedOption = event.target.value as partitionOffsetOption;
        let startOffset = defaultStartOffset;
        let endOffset = defaultEndOffset;

        if (selectedOption === 'Latest') {
            startOffset = defaultEndOffset - numMessages + 1;
            endOffset = defaultEndOffset;
        } else if (selectedOption === 'Earliest') {
            startOffset = 0;
            endOffset = numMessages - 1;
        } else if (selectedOption === 'Custom') {
            // Custom logic can be added here if needed
            // For now, we will use the default values
        }

        onPartitionDetailsChange(startOffset, endOffset);
        setOffsets({ startOffset, endOffset });
    }

    const [offsets, setOffsets] = React.useState<offsetState>(
        { startOffset: defaultStartOffset, endOffset: defaultEndOffset }
    );
    const [numMessages, setNumMessages] = React.useState<number>(defaultNumMessages);

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

            <div
                style={{
                    display: 'flex',
                    gap: '0.5rem',
                    alignItems: 'center',
                    flexDirection: 'column',
                    fontSize: '0.875rem',
                }}
            >
                Offset Range
                <select defaultValue = 'Latest' onChange={handleOffsetChange}>
                    {
                        partitionOffsetOptions.map((option) => (
                            <option key={option} value={option}>
                                {option}
                            </option>
                        ))
                    }
                </select>
            </div>
        </div>
    )
}

export default TopicDetailsOptionBar;