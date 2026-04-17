import React from 'react';
import { render, screen } from '@testing-library/react';
import { ProgressBar } from '@/components/ui/ProgressBar';
import { Breadcrumb } from '@/components/ui/Breadcrumb';
import { Button } from '@/components/ui/Button';
import '@testing-library/jest-dom';

describe('UI Components', () => {
  describe('ProgressBar', () => {
    it('renders the progress bar container', () => {
      const { container } = render(<ProgressBar progress={50} />);
      const innerBar = container.querySelector('.bg-slate-100');
      expect(innerBar).toBeInTheDocument();
    });
  });

  describe('Breadcrumb', () => {
    it('renders the correct labels and links', () => {
      render(<Breadcrumb items={[
        { label: 'Home', href: '/' },
        { label: 'Current' }
      ]} />);
      
      expect(screen.getByText('Home')).toBeInTheDocument();
      expect(screen.getByText('Current')).toBeInTheDocument();
      
      const homeLink = screen.getByText('Home');
      expect(homeLink.closest('a')).toHaveAttribute('href', '/');
    });
  });

  describe('Button', () => {
    it('renders children correctly', () => {
      render(<Button>Click Me</Button>);
      expect(screen.getByText('Click Me')).toBeInTheDocument();
    });

    it('renders likert active state correctly', () => {
      render(<Button variant="likert" isActive={true}>Likert</Button>);
      const btn = screen.getByText('Likert').closest('button');
      expect(btn?.className).toContain('border-palette4');
    });
  });
});
