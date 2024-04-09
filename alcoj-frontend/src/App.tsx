/* eslint-disable react/no-array-index-key */
import React from 'react';
import { Layout, Input, Select, Button, Card } from 'antd';
import { CaretRightOutlined } from '@ant-design/icons';
import { useState } from 'react';
import axios from 'axios';

const { Header, Footer, Sider, Content } = Layout;

const backgroundColor = '#F0F0F0';

const headerStyle: React.CSSProperties = {
  textAlign: 'left',
  color: '#000',
  height: 64,
  paddingInline: 48,
  lineHeight: '64px',
  backgroundColor: backgroundColor,
  fontSize: 24,
};

const siderStyle: React.CSSProperties = {
  textAlign: 'center',
  lineHeight: '120px',
  color: '#fff',
  backgroundColor: backgroundColor,
  borderRadius: 10,
  display: 'flex',
  flexDirection: 'column',
};

const footerStyle: React.CSSProperties = {
  display: 'flex',
  color: '#fff',
  backgroundColor: backgroundColor,
  height: '25vh',
  padding: '0',
};

const layoutStyle = {
  borderRadius: 8,
  overflow: 'hidden',
  width: '100vw',
  height: '100vh',
};

const codeSpaceStyle: React.CSSProperties = {
  width: '100%',
  height: '100%',
  backgroundColor: '#FFF',
  padding: '10px',
  resize: 'none',
  border: 'none',
  boxShadow: 'none',

  lineHeight: 1.5,
  fontSize: 16,
};

const contentStyle: React.CSSProperties = {
  display: 'flex',
  textAlign: 'left',
  minHeight: 120,
  lineHeight: '120px',
  borderRadius: 10,
  backgroundColor: backgroundColor,
  padding: '10px 5px 10px 10px',
};

const handleChange = (value: string) => {
  console.log(`selected ${value}`);
};

const leftStyle: React.CSSProperties = {
  display: 'flex',
  flexDirection: 'column',
  gap: '5px',
  textAlign: 'left',
  borderRadius: 10,
  width: '100%',
  position: 'relative',
  backgroundColor: '#FFF',
};

const runButtonStyle: React.CSSProperties = {
  width: '50px',
  height: '50px',
  borderRadius: 25,
  textAlign: 'center',
  position: 'absolute',
  right: 10,
  bottom: 10,
};

const rightStyle: React.CSSProperties = {
  display: 'flex',
  flexDirection: 'column',
  gap: '10px',
  textAlign: 'left',
  borderRadius: 10,
  width: '100%',
  height: '100%',
  alignItems: 'center',
  padding: '10px',
};

const bottomStyle: React.CSSProperties = {
  display: 'flex',
  flexDirection: 'row',
  height: '100%',
  width: '100%',
  padding: '0 10px',
  gap: '10px',
  backgroundColor: backgroundColor,
}

interface RunResponse {
  output: string;
  static_analysis: StaticAnalysis[];
  runtime_analysis: Record<string, string>;
}

interface StaticAnalysis {
  row: string;
  column: string;
  message: string;
}

const App: React.FC = () => {
  const [selectedLanguage, setSelectedLanguage] = useState<string>('python');
  const [code, setCode] = useState<string>('');
  const [input, setInput] = useState<string>('');

  const [response, setResponse] = useState<RunResponse | null>(null);
  let formattedOutput = response?.output?.replace(/\r?\n/g, '<br>');

  const handleChange = (value: string) => {
    setSelectedLanguage(value);
  };

  const handleCodeChange = (event: React.ChangeEvent<HTMLTextAreaElement>) => {
    setCode(event.target.value);
  };

  const handleInputChange = (event: React.ChangeEvent<HTMLTextAreaElement>) => {
    setInput(event.target.value);
  };

  const handleRunButtonClick = async () => {
    const requestBody = {
      language: selectedLanguage,
      code: code,
      input: input,
    };
    console.log('Request body:', requestBody);

    try {
      let url = 'http://localhost:8080/alcoj/api/v1';
      switch (selectedLanguage) {
        case 'python':
          url = 'http://localhost:8080/alcoj/api/v1';
          break;
        case 'golang':
          url = 'http://localhost:8081/alcoj/api/v1';
          break;
        default:
          break;
      }
      const response = await axios.post(url, requestBody);

      if (response.status === 200 && response.data) {
        setResponse(response.data);
        console.log('Response:', response.data);
      }
    } catch (error) {
      console.error('Error fetching data:', error);
    }
  };

  return (
    <Layout style={layoutStyle}>
      <Header style={headerStyle}>
        AlcOJ
      </Header>
      <Layout style={{ color: "#FFF", height: '60vh' }}>
        <Content style={contentStyle}>
          <div style={leftStyle}>
            <Select
              defaultValue="python"
              style={{ width: 120 }}
              onChange={handleChange}
              options={[
                { value: 'python', label: 'Python3' },
                { value: 'golang', label: 'Golang' },
              ]}
            />
            <Input.TextArea
              style={{ ...codeSpaceStyle, }}
              placeholder="Code Space"
              onChange={handleCodeChange}
            />
            <Button type="primary"
              style={runButtonStyle}
              icon={<CaretRightOutlined />}
              onClick={handleRunButtonClick}
            />
          </div>
        </Content>
        <Sider width="25%" style={siderStyle}>
          <div style={rightStyle}>
            <Card style={{ width: "100%", height: "50%", overflowY: 'auto', }} title='Static Analysis'>
              {response &&
                response.static_analysis &&
                response.static_analysis.map((analysis, index) => (
                  <div key={index}>
                    <p>
                      {analysis.row}:{analysis.column} {analysis.message}
                    </p>
                  </div>
                ))}
            </Card>
            <Card style={{ width: "100%", height: "50%", overflowY: 'auto', }} title='Runtime Analysis'>
              {response && Object.keys(response.runtime_analysis).map((key, index) => (
                <div key={index}>
                  <p>{key}: {response.runtime_analysis[key]}</p>
                </div>
              ))}
            </Card>
          </div>
        </Sider>
      </Layout>

      <Footer style={footerStyle}>
        <div style={bottomStyle}>
          <Input.TextArea
            style={{ ...codeSpaceStyle, flex: 1 }}
            placeholder="Input Space"
            onChange={handleInputChange}
          />
          <Card style={{ flex: 1, overflowY: 'auto', }} title="Outputs">
            {formattedOutput ? <div dangerouslySetInnerHTML={{ __html: formattedOutput }} /> : null}
          </Card>
        </div>
      </Footer>
    </Layout>
  );
};

export default App;
