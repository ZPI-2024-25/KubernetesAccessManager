// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-expect-error
import yaml from 'js-yaml';

/**
 * Converts a JSON object to a YAML string.
 * @param json - The JSON object to convert.
 * @returns The YAML string.
 */
export const jsonToYaml = (json: unknown): string => {
    try {
        return yaml.dump(json);
    } catch (error) {
        console.error("Error converting JSON to YAML:", error);
        throw new Error("Failed to convert JSON to YAML");
    }
};

/**
 * Converts a YAML string to a JSON object.
 * @param yamlString - The YAML string to convert.
 * @returns The JSON object.
 */
export const yamlToJson = (yamlString: string): unknown => {
    try {
        return yaml.load(yamlString);
    } catch (error) {
        console.error("Error converting YAML to JSON:", error);
        throw new Error("Failed to convert YAML to JSON");
    }
};

/**
 * Parses a JSON string into a specified type.
 * @param jsonString - The JSON string to parse.
 * @returns The parsed object of type T.
 */
export const parseJson = <T>(jsonString: string): T => {
    try {
        return JSON.parse(jsonString) as T;
    } catch (error) {
        console.error("Error parsing JSON:", error);
        throw new Error("Failed to parse JSON");
    }
};

/**
 * Converts a JSON object to a formatted JSON string.
 * @param json - The JSON object to convert.
 * @returns The formatted JSON string.
 */
export const stringifyJson = (json: unknown): string => {
    try {
        return JSON.stringify(json, null, 2); // Formatted JSON
    } catch (error) {
        console.error("Error converting JSON to string:", error);
        throw new Error("Failed to convert JSON to string");
    }
};

/**
 * Parses a YAML string into a specified type.
 * @param yamlString - The YAML string to parse.
 * @returns The parsed object of type T.
 */
export const parseYaml = <T>(yamlString: string): T => {
    try {
        return yaml.load(yamlString) as T;
    } catch (error) {
        console.error("Error parsing YAML:", error);
        throw new Error("Failed to parse YAML");
    }
};

/**
 * Converts a JSON object to a YAML string.
 * @param json - The JSON object to convert.
 * @returns The YAML string.
 */
export const stringifyYaml = (json: unknown): string => {
    try {
        return yaml.dump(json);
    } catch (error) {
        console.error("Error converting JSON to YAML:", error);
        throw new Error("Failed to convert JSON to YAML");
    }
};