import { extendTheme } from "@mui/joy/styles";

const theme = extendTheme({
  fontFamily: {
    body: "'Public Sans', var(--joy-fontFamily-fallback)",
  },
});

export default theme;
