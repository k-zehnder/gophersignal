import Typography from '@mui/joy/Typography';
import Box from '@mui/joy/Box';
import Link from '@mui/joy/Link';

function Description() {
  return (
    <Box sx={{ my: 4 }}>
      <Typography level="h4" component="h3" sx={{ fontWeight: 'medium' }}>
        Welcome to Gopher Signal
      </Typography>
      <Typography sx={{ mt: 1 }}>
        Gopher Signal uses ChatGPT to quickly summarize important points from{" "}
        <Link href="https://news.ycombinator.com" target="_blank" rel="noopener noreferrer">
          Hacker News
        </Link>{" "}
        articles, giving you brief and useful updates.
      </Typography>
    </Box>
  );
}

export default Description;
