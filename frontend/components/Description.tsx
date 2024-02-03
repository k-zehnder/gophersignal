/**
 * Description: A component that renders a welcome message and a brief explanation
 * of the Gopher Signal application. It provides an introductory text and a link to
 * Hacker News, giving context to the application's purpose and functionality.
 */
import Typography from '@mui/joy/Typography';
import Box from '@mui/joy/Box';
import Link from '@mui/joy/Link';

// Description: A component that renders a welcome message and a brief explanation
function Description() {
  return (
    <Box sx={{ my: 4 }}>
      {/* Display a heading with a welcome message */}
      <Typography level="h4" component="h3" sx={{ fontWeight: 'medium' }}>
        Welcome to Gopher Signal
      </Typography>

      {/* Display a paragraph with an explanation of Gopher Signal */}
      <Typography sx={{ mt: 1 }}>
        Gopher Signal uses smart technology to quickly summarize important points from{" "}
        {/* Create a link to Hacker News */}
        <Link href="https://news.ycombinator.com" target="_blank" rel="noopener noreferrer">
          Hacker News
        </Link>{" "}
        articles, giving you brief and useful updates.
      </Typography>
    </Box>
  );
}

export default Description;
