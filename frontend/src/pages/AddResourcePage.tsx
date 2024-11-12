import React from 'react';
import { InputForm } from '../components/AddForm/InputForm';  // Форма для добавления ресурса

export const AddResourcePage: React.FC = () => {
    return (
        <div>
            <h2>Add Resource</h2>
            <InputForm />
        </div>
    );
};

