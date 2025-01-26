import { useEffect } from 'react';
import { useColorScheme } from '@mui/joy/styles';
import IconButton from '@mui/joy/IconButton';
import * as React from 'react';
import LightModeIcon from '@mui/icons-material/LightMode';
import DarkModeIcon from '@mui/icons-material/DarkMode';

// ModeButton component allows users to toggle between light and dark modes.
export default function ModeButton() {
  const { mode, setMode } = useColorScheme();
  const [mounted, setMounted] = React.useState<boolean>(false);

  useEffect(() => {
    // Set `mounted` to `true` when the component is mounted.
    setMounted(true);
  }, []);

  if (!mounted) {
    // Render a placeholder button to avoid layout shift before mounting.
    return (
      <IconButton
        variant="plain"
        color="neutral"
        sx={{ width: 60 }}
        aria-label="Toggle light and dark mode"
      />
    );
  }

  // Render the mode toggle button based on the current mode.
  return (
    <IconButton
      variant="plain"
      color="neutral"
      aria-label="Toggle light and dark mode"
      onClick={() => setMode(mode === 'dark' ? 'light' : 'dark')}
    >
      {mode === 'dark' ? <LightModeIcon /> : <DarkModeIcon />}
    </IconButton>
  );
}
