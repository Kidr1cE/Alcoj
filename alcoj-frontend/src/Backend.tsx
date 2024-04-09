/* eslint-disable react/no-array-index-key */
import React, { useState, useRef, useLayoutEffect, useEffect } from 'react';
import { Layout, Select, Button, Card, Typography, Flex, Col, Row, Statistic } from 'antd';
import CountUp from 'react-countup';
import axios from 'axios'; // 假设您已经安装了axios库来处理HTTP请求

const { Header, Content } = Layout;
const { Paragraph, Text } = Typography;

const backgroundColor = '#F0F0F0';

type Worker = {
    worker_id: string;
    worker_status: string;
}

type Request = {
    key: string;
}

type Response = {
    worker_num: number;
    queue_tasks: number;
    finished_tasks: number;
    workers: Worker[];
}

const WorkerCard: React.FC<Worker> = ({ worker_id, worker_status }) => {
    const truncatedWorkerId = worker_id.slice(0, 8);

    const getStatusColor = (status: string) => {
        switch (status) {
            case 'ready':
                return '#73d13d';
            case 'running':
                return '#ffc53d';
            case 'error':
                return '#ff4d4f';
            default:
                return '#4096ff';
        }
    };

    return (
        <Card bordered={false} style={{
            paddingBottom: '1px',
            textAlign: 'center',
            margin: '5px',
        }}>
            <div style={{
                width: '100px',
                height: '100px',
                backgroundColor: getStatusColor(worker_status),
                borderRadius: '8px',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
            }}>
                {worker_status}
            </div>
            <Paragraph copyable={{ text: worker_id }}>{truncatedWorkerId}</Paragraph>
        </Card>
    );
};

const WorkersList: React.FC<{ workers: Worker[] }> = ({ workers }) => (
    <Flex>
        {workers.map((worker) => (
            <WorkerCard key={worker.worker_id} {...worker} />
        ))}
    </Flex>
);


const ResponseDisplay: React.FC<Response> = ({ worker_num, queue_tasks, finished_tasks, workers }) => {
    const formatter = (value: any) => <CountUp
        start={0}
        end={value}
        useEasing={false}
        separator=","
    />;
    return (
        <div style={{
            backgroundColor: '#fff',
            margin: '10px',
            borderRadius: '8px',
            padding: '10px',
            height: 'calc(100vh - 120px)',
        }}>
            <Row gutter={21}>
                <Col span={7}>
                    <Statistic title="Online Workers" value={worker_num} formatter={formatter} />
                </Col>
                <Col span={7}>
                    <Statistic title="Queue Tasks" value={queue_tasks} />
                </Col>
                <Col span={7}>
                    <Statistic title="Finished Tasks" value={finished_tasks} />
                </Col>
            </Row>
            <WorkersList workers={workers} />
        </div>
    );
};
const initialParsedResponse: Response = {
    worker_num: 6,
    queue_tasks: 20,
    finished_tasks: 114514,
    workers: [
        { worker_id: 'acf9b8ba-f014-11ee-a043-00155da91e79', worker_status: 'Running' },
        { worker_id: '9046dff8-f015-11ee-a604-00155da91e79', worker_status: 'Ready' },
    ],
};

const Backend: React.FC = () => {
    const wsRef = useRef<WebSocket | null>(null);
    const [selectedLanguage, setSelectedLanguage] = useState<string>('python');
    const [message, setMessage] = useState('');
    const [parsedResponse, setParsedResponse] = useState<Response>(initialParsedResponse);


    const handleChange = (value: string) => {
        setSelectedLanguage(value);
    };

    useEffect(() => {
        let url = `ws://localhost:7070/${selectedLanguage}`;
        switch (selectedLanguage) {
            case 'python':
                url = 'ws://localhost:7070/python';
                break;
            case 'golang':
                url = 'ws://localhost:7071/golang';
                break;
            default:
                break;
        }
        wsRef.current = new WebSocket(url);
        wsRef.current.onmessage = (e) => {
            setMessage(e.data);
            try {
                const response: Response = JSON.parse(e.data);
                setParsedResponse(response);
            } catch (error) {
                console.error('Failed to parse response:', error);
            }
        };

        return () => {
            wsRef.current?.close();
        };
    }, [selectedLanguage]);

    useEffect(() => {
        let intervalId: NodeJS.Timeout | null = null;

        const sendRequest = () => {
            if (wsRef.current) {
                const request: Request = { key: 'ssss' };
                wsRef.current.send(JSON.stringify(request));
            }
        };

        intervalId = setInterval(sendRequest, 1000);

        return () => {
            clearInterval(intervalId!);
        };
    }, [wsRef.current]);

    return (
        <Layout style={{
            borderRadius: 8,
            overflow: 'hidden',
            width: '100vw',
            height: '100vh',
        }}>
            <Header style={{
                textAlign: 'left',
                color: '#000',
                height: 64,
                paddingInline: 48,
                lineHeight: '64px',
                backgroundColor: backgroundColor,
                fontSize: 24,
            }}>
                AlcOJ Backend Preview
                <Select
                    defaultValue="python3"
                    style={{ width: 120, marginLeft: 20 }}
                    onChange={handleChange}
                    options={[
                        { value: 'python', label: 'Python3' },
                        { value: 'golang', label: 'Golang' },
                    ]}
                />
            </Header>
            <Content style={{
                textAlign: 'center',
                lineHeight: '120px',
                color: '#fff',
                backgroundColor: '#F0F0F0',
            }}>
                <ResponseDisplay {...parsedResponse} />
            </Content>
        </Layout>
    );
};

export default Backend;
