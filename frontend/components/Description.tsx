import Typography from '@mui/joy/Typography';
import Box from '@mui/joy/Box';
import Link from '@mui/joy/Link';
function Description() {
  return (
    <Box sx={{ my: 4 }}> 
      <Typography level="h4" component="h3" sx={{ fontWeight: 'medium' }}>
        Welcome to Gopher Signal
      </Typography>
      <Typography level="body1" sx={{ mt: 1 }}>
        Your concise source of summarized tech news. We use ChatGPT to provide quick summaries of the latest scraped articles from 
        <Link href="https://news.ycombinator.com" target="_blank" rel="noopener noreferrer" sx={{ ml: 0.5 }}>
          Hacker News
        </Link>.
      </Typography>
    </Box>
  );
}

export default Description;
