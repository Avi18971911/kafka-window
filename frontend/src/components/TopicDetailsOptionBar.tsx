import React, {useEffect, useRef} from "react";

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
    const [offsets, setOffsets] = React.useState<offsetState>(
        { startOffset: defaultStartOffset, endOffset: defaultEndOffset }
    );
    const [numMessages, setNumMessages] = React.useState<number>(defaultNumMessages);
    const [offsetCategory, setOffsetCategory] = React.useState<partitionOffsetOption>('Latest');
    const isCustomOffset = offsetCategory === 'Custom';

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
    const debounceTimeoutRef = useRef<number | null>(null);

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
            startOffset = offsets.startOffset
            endOffset = offsets.endOffset;
        }

        setOffsetCategory(selectedOption);
        onPartitionDetailsChange(startOffset, endOffset);
        setOffsets({ startOffset, endOffset });
    }

    const updatePartitionDetailsFromMessagesChange = (value: number) => {
        let startOffset = -1;
        let endOffset = -1;

        switch (offsetCategory) {
            case 'Latest':
                startOffset = -1 * value;
                endOffset = -1;
                setOffsets({
                    startOffset: startOffset,
                    endOffset: endOffset
                });
                onPartitionDetailsChange(startOffset, endOffset);
                break;
            case 'Earliest':
                startOffset = 0;
                endOffset = value - 1;
                setOffsets({
                    startOffset: startOffset,
                    endOffset: endOffset
                });
                onPartitionDetailsChange(startOffset, endOffset);
                break;
            case 'Custom':
                // Custom will hide the num messages change logic
                break;
            default:
                console.error("Unknown offset category:", offsetCategory);
        }
    };


    const handleNumMessagesChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        const value = parseInt(event.target.value, 10);
        if (isNaN(value) || value <= 0) {
            console.error("Invalid number of messages:", event.target.value);
            return;
        }
        setNumMessages(value);

        if (debounceTimeoutRef.current) {
            clearTimeout(debounceTimeoutRef.current);
        }

        debounceTimeoutRef.current = setTimeout(() => {
            updatePartitionDetailsFromMessagesChange(value);
        }, 500);
    }

    // Cleanup debounce timeout on unmount
    useEffect(() => {
        return () => {
            if (debounceTimeoutRef.current) {
                clearTimeout(debounceTimeoutRef.current);
            }
        };
    }, []);

    const handleStartOffsetChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        const value = parseInt(event.target.value, 10);
        if (isNaN(value) || value == 0) {
            console.error("Invalid start offset:", event.target.value);
            return;
        }
        setOffsets((prev) => ({ ...prev, startOffset: value }));
        onPartitionDetailsChange(value, offsets.endOffset);
    }

    const handleEndOffsetChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        const value = parseInt(event.target.value, 10);
        if (isNaN(value) || value == 0) {
            console.error("Invalid end offset:", event.target.value);
            return;
        }
        setOffsets((prev) => ({ ...prev, endOffset: value }));
        onPartitionDetailsChange(offsets.startOffset, value);
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

            {
                isCustomOffset ?
                    <div
                        style={{
                            display: 'flex',
                            gap: '1.0rem',
                            alignItems: 'center',
                            flexDirection: 'row',
                            fontSize: '0.875rem',
                        }}
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
                            Start Offset
                            <input
                                type='number'
                                value={offsets.startOffset}
                                onChange={handleStartOffsetChange}
                                style={{ width: '80px' }}
                            />
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
                            End Offset
                            <input
                                type='number'
                                value={offsets.endOffset}
                                onChange={handleEndOffsetChange}
                                style={{ width: '80px' }}
                            />
                        </div>
                    </div>
                :
                    <div
                        style={{
                            display: 'flex',
                            gap: '0.5rem',
                            alignItems: 'center',
                            flexDirection: 'column',
                            fontSize: '0.875rem',
                        }}
                    >
                        Number of Messages
                        <input
                            type='number'
                            value={numMessages}
                            onChange={handleNumMessagesChange}
                            style={{ width: '80px' }}
                        />
                    </div>
            }
        </div>
    )
}

export default TopicDetailsOptionBar;