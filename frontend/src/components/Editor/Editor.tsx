import {Button, Card, message} from "antd";
import style from "./Editor.module.css";
import {Editor as MonacoEditor} from "@monaco-editor/react";
import {useCallback, useEffect, useState} from "react";
import LanguageSelector from "./LanguageSelector.tsx";
import {stringifyJson, parseJson, parseYaml, stringifyYaml} from "../../functions/jsonYamlFunctions.ts";
import {ResourceDetails} from "../../types";
import {useNavigate} from "react-router-dom";

const Editor = ({name, text, endpoint, namespaceSelector}: {
    name: string,
    text: string,
    endpoint: (data: ResourceDetails) => Promise<ResourceDetails>
    namespaceSelector?: React.ReactNode
}) => {
    const [value, setValue] = useState<string>(text);
    const [language, setLanguage] = useState<string>("yaml");
    const navigate = useNavigate();

    useEffect(() => {
        setValue(text);
    }, [text]);

    const onLanguageChange = useCallback(
        (newLanguage: string) => {
            let convertedValue = value;

            try {
                if (newLanguage === 'json' && language === 'yaml') {
                    if (value.trim() !== '') {
                        const jsonValue = parseYaml(value);
                        convertedValue = stringifyJson(jsonValue);
                    } else {
                        convertedValue = '';
                    }
                } else if (newLanguage === 'yaml' && language === 'json') {
                    if (value.trim() !== '') {
                        const jsonValue = parseJson(value);
                        convertedValue = stringifyYaml(jsonValue);
                    } else {
                        convertedValue = '';
                    }
                }
                setLanguage(newLanguage);
                setValue(convertedValue);
            } catch (error) {
                console.error('Conversion error:', error);
                message.error('Conversion failed: Check the input format.');
            }
        },
        [language, value]
    );

    const isEmptyOrWhitespace = (str: string): boolean => {
        return !str || str.trim() === '';
    };

    const onSave = useCallback(async () => {
        if (isEmptyOrWhitespace(value)) {
            message.warning('Cannot save empty content.');
            return;
        }

        try {
            let parsedData: ResourceDetails;
            if (language === 'json') {
                parsedData = parseJson<ResourceDetails>(value);
            } else {
                parsedData = parseYaml<ResourceDetails>(value);
            }
            console.log('Parsed Data:', parsedData);

            try {
                const response = await endpoint(parsedData);
                console.log('Save successful:', response);
                message.success('Save successful');
                // setValue(stringifyYaml(response));
                navigate(-1);
            } catch (error) {
                if (error instanceof Error) {
                    console.error('Error getting resource:', error);
                    message.error(error.message, 4);
                } else {
                    message.error('An unexpected error occurred.');
                }
            }
        } catch (error) {
            console.error('Invalid format', error);
            message.error('Save failed: Invalid format');
        }
    }, [language, value, endpoint]);

    const handleEditorChange = useCallback((newValue: string | undefined) => {
        setValue(newValue || '');
    }, []);

    return (
        <Card className={style.content} style={{marginTop: '64px',}} title={name}>
            <div className={style.editorOptionsPanel}>
                <div className={style.selectors}>
                    <LanguageSelector language={language} onSelect={onLanguageChange}/>

                    {namespaceSelector}
                </div>

                <div style={{display: 'flex', gap: '8px'}}>
                    <Button type="default" onClick={() => navigate(-1)}>Back</Button>
                    <Button type="primary" onClick={onSave}>Save</Button>
                </div>
            </div>
            <MonacoEditor
                height="65vh"
                width="80wh"
                theme="vs-dark"
                language={language}
                value={value}
                onChange={handleEditorChange}
            />
        </Card>
    );
};

export default Editor;