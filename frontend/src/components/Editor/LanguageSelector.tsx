import {Select} from "antd";

const LanguageSelector = ({language, onSelect}: { language: string, onSelect: (value: string) => void }) => {
    return (
        <Select
            value={language}
            style={{width: 120}}
            onChange={onSelect}
            options={[
                {value: 'yaml'},
                {value: 'json'},
            ]}
        />
    );
};

export default LanguageSelector;