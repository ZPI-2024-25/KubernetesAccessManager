import {Button, Card, message} from "antd";
import style from "./Editor.module.css";
import {Editor as MonacoEditor} from "@monaco-editor/react";
import * as monaco from "monaco-editor";
import {useCallback, useEffect, useRef, useState} from "react";
import LanguageSelector from "./LanguageSelector.tsx";
import {stringifyJson, parseJson, parseYaml, stringifyYaml} from "../../functions/jsonYamlFunctions.ts";
import {ResourceDetails} from "../../types/ResourceDetails.ts";

const Editor = ({name, text, endpoint}: {
    name: string,
    text: string,
    endpoint: (data: ResourceDetails) => Promise<ResourceDetails>
}) => {
    const editorRef = useRef<monaco.editor.IStandaloneCodeEditor | null>(null);
    const [value, setValue] = useState<string>(text);
    const [language, setLanguage] = useState<string>("yaml");

    // TODO: zamienić wzorzec na stworzenie
    useEffect(() => {
        console.log(text);
        //setValue(text);
        setValue(stringifyYaml((JSON.stringify(text, null, 2))));
    }, [text]);

    const onMount = useCallback((editor: monaco.editor.IStandaloneCodeEditor) => {
        editorRef.current = editor;
        editor.focus();
    }, []);

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
            } catch (error) {
                console.error('Save failed:', error);
                message.error('Save failed');
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
        <Card className={style.content} title={name}>
            <div className={style.editorOptionsPanel}>
                <LanguageSelector language={language} onSelect={onLanguageChange}/>
                <Button type="primary" onClick={onSave}>Save</Button>
            </div>
            <MonacoEditor
                height="65vh"
                width="80wh"
                theme="vs-dark"
                language={language}
                onMount={onMount}
                value={value}
                onChange={handleEditorChange}
            />
        </Card>
    );
};

export default Editor;