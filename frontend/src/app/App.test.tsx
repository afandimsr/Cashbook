import { render, screen } from '@testing-library/react';
import App from './App';
import { describe, it, expect } from 'vitest';

describe('App', () => {
    it('renders dashboard by default', () => {
        render(<App />);
        expect(screen.getByText(/Login System/i)).toBeInTheDocument();
    });
});
