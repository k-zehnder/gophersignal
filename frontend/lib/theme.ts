/**
 * theme: Customizes the default Material-UI theme.
 * This theme configuration creates a custom theme by extending the default Material-UI theme.
 * It defines custom font family settings for the 'body' text, ensuring a consistent typography style.
*/
import { extendTheme } from "@mui/joy/styles";

// Create a custom theme by extending the default Material-UI theme.
const theme = extendTheme({
  // Define custom font family settings for the 'body' text.
  fontFamily: {
    body: "'Public Sans', var(--joy-fontFamily-fallback)",
  },
});

// Export the custom theme to make it available for use in the application.
export default theme;
