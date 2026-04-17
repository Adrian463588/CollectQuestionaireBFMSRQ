import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import IPIPQuestionnaire from '../app/questionnaire/ipip/page';
import { ipipBfm50Questions } from '@/data/ipip';
import '@testing-library/jest-dom';

describe('IPIPQuestionnaire Interactions', () => {
  it('preserves the active state of an answer when returning to a previous question', async () => {
    render(<IPIPQuestionnaire />);

    expect(screen.getByText(ipipBfm50Questions[0].text)).toBeInTheDocument();

    const btnSangatSesuai = screen.getByText('5').closest('button');
    expect(btnSangatSesuai).not.toBeNull();
    fireEvent.click(btnSangatSesuai!);

    // Should navigate to question 2
    expect(await screen.findByText(ipipBfm50Questions[1].text)).toBeInTheDocument();

    // Click "Kembali"
    const btnKembali = screen.getByText('Kembali');
    fireEvent.click(btnKembali);

    // Validate we are back on question 1
    expect(await screen.findByText(ipipBfm50Questions[0].text)).toBeInTheDocument();

    // Check if the Sangat Sesuai button is highlighted (active likert)
    const btnSSActive = screen.getByText('5').closest('button');
    expect(btnSSActive?.className).toContain('border-palette4');
  });
});
