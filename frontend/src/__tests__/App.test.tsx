//https://dev.to/hannahadora/jest-testing-with-vite-and-react-typescript-4bap

import '@testing-library/jest-dom'
import { render } from "@testing-library/react"
import App from "../App.tsx";

test('demo', () => {
    expect(1).toBe(1);
})

test("Renders the main page", () => {
    render(<App />)
    expect(true).toBeTruthy()
})