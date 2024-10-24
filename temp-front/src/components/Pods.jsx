import { useEffect, useState } from "react";

function Pods({ token }) {
    const [pods, setPods] = useState([]);
    const [columns, setColumns] = useState([]);
    const [error, setError] = useState(null);

    useEffect(() => {
        const fetchPods = async () => {
            try {
                const response = await fetch("http://localhost:8080/api/v1/k8s/Pod", {
                    method: "GET",
                    headers: {
                        Authorization: `Bearer ${token}`,
                        "Content-Type": "application/json",
                    },
                });

                if (!response.ok) {
                    throw new Error(`Error: ${response.status}`);
                }

                const data = await response.json();
                setColumns(data.columns);
                setPods(data.resource_list);
            } catch (err) {
                setError(err.message);
            }
        };

        fetchPods();
    }, [token]);

    if (error) {
        return <div>Error: {error}</div>;
    }

    const renderCell = (pod, column) => {
        return pod[column] || "N/A";
    };

    return (
        <div>
            <h2>Pods List</h2>
            <table>
                <thead>
                <tr>
                    {columns.map((column, index) => (
                        <th key={index}>{column}</th>
                    ))}
                </tr>
                </thead>
                <tbody>
                {pods.map((pod, index) => (
                    <tr key={index}>
                        {columns.map((column, colIndex) => (
                            <td key={colIndex}>{renderCell(pod, column)}</td>
                        ))}
                    </tr>
                ))}
                </tbody>
            </table>
        </div>
    );
}

export default Pods;
