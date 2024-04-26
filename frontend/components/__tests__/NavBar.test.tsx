import React from 'react';
import '@testing-library/jest-dom';
import { render, screen } from '@testing-library/react';
import NavBar from '../../components/NavBar';
import { CssVarsProvider } from '@mui/joy/styles';
import { siteMetaData } from '../../lib/siteMetaData';

// Mocking window.matchMedia to ensure it can be used in tests as it's not implemented in jsdom
window.matchMedia = jest.fn().mockImplementation((query) => ({
  matches: false, // Default matches to false unless specified otherwise
  media: query,
  onchange: null,
  addListener: jest.fn(), // Mock implementation for addListener
  removeListener: jest.fn(), // Mock implementation for removeListener
}));

describe('NavBar Component', () => {
  it('renders the navigation links and the site title', () => {
    render(
      <CssVarsProvider>
        <NavBar />
      </CssVarsProvider>,
    );

    // Check for the site title as a heading to verify it is rendered correctly
    const siteTitle = screen.getByRole('heading', { name: siteMetaData.title });
    expect(siteTitle).toBeInTheDocument();

    // Check that navigation links are rendered and properly labeled
    const homeLink = screen.getByRole('link', { name: 'Home' });
    expect(homeLink).toBeInTheDocument();
    expect(homeLink).toHaveAttribute('href', '/');

    const aboutLink = screen.getByRole('link', { name: 'About' });
    expect(aboutLink).toBeInTheDocument();
    expect(aboutLink).toHaveAttribute('href', '/about');

    // Check if the API link is correct based on environment settings
    const apiUrl =
      process.env.NEXT_PUBLIC_ENV === 'development'
        ? 'http://localhost:8080/swagger/index.html#/'
        : 'https://gophersignal.com/swagger/index.html#/';
    const apiLink = screen.getByRole('link', { name: 'API' });
    expect(apiLink).toBeInTheDocument();
    expect(apiLink).toHaveAttribute('href', apiUrl);

    // Locate the ModeButton by its aria-label to ensure it's accessible and functioning
    const modeButton = screen.getByLabelText('Toggle light and dark mode');
    expect(modeButton).toBeInTheDocument();
  });
});
