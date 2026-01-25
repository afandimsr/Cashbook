import { render, screen } from '@testing-library/react';
import App from './App';
import { describe, it, expect } from 'vitest';

describe('App', () => {
    it('renders dashboard by default', () => {
        render(<App initialTitle="Login System" />);
        expect(screen.getByText(/Login System/i)).toBeInTheDocument();
    });
});
