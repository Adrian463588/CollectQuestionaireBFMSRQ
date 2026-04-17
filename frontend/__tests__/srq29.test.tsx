import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import SRQ29Questionnaire from '../app/questionnaire/srq29/page';
import { srq29Questions } from '@/data/srq29';
import '@testing-library/jest-dom';

// Note: Mocks are already configured in jest.setup.ts
// including next/navigation and @/lib/api

describe('SRQ29Questionnaire Bug Fix Validation', () => {
  it('preserves the active state of an answer when returning to a previous question', async () => {
    render(<SRQ29Questionnaire />);

    // Validate we are on the first question
    expect(screen.getByText(srq29Questions[0].text)).toBeInTheDocument();

    // Click "Ya" for the first question
    const btnYa = screen.getByText('Ya').closest('button');
    expect(btnYa).not.toBeNull();
    fireEvent.click(btnYa!);

    // Clicking 'Ya' should navigate to the second question eventually (await due to possible animation/state batching)
    expect(await screen.findByText(srq29Questions[1].text)).toBeInTheDocument();

    // Click "Kembali" to return to the first question
    const btnKembali = screen.getByText('Kembali');
    fireEvent.click(btnKembali);

    // Validate we are back on the first question
    expect(await screen.findByText(srq29Questions[0].text)).toBeInTheDocument();

    // Re-select the buttons on the screen since DOM elements might be re-rendered
    const btnYaActive = screen.getByText('Ya').closest('button');
    const btnTidakInactive = screen.getByText('Tidak').closest('button');
    
    expect(btnYaActive?.className).toContain('bg-palette4/10');
    expect(btnYaActive?.className).toContain('border-palette4');
    expect(btnYaActive?.className).toContain('text-palette4');

    expect(btnTidakInactive?.className).not.toContain('bg-palette4/10');
  });
});
