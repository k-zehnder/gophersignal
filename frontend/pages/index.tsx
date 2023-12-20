import Layout from "../components/Layout";
import Link from "next/link";
import Date from "../components/Date";
import { GetStaticProps } from "next";
import Typography from "@mui/joy/Typography";
import List from "@mui/joy/List";
import ListItem from "@mui/joy/ListItem";
import { getSortedPostsData } from "../lib/posts";

export default function Index({
  allPostsData,
}: {
  allPostsData: {
    title: string;
    summary: string;
    category: string;
    date: string;
    id: string;
  }[];
}) {
  return (
    <Layout>
      {/* Header */}
      <Typography level="h1" component="h2" sx={{ fontWeight: 'bold', mb: 2 }}>
        Latest Posts
      </Typography>

      {/* Posts List */}
      <List sx={{ display: "flex", flexDirection: 'column', gap: "1rem" }}>
        {allPostsData.map(({ id, date, category, title, summary }) => (
          <ListItem
            key={id}
            sx={{
              display: "flex",
              flexDirection: "column",
              alignItems: "flex-start",
              '&:hover a': {
                textDecoration: 'underline',
              },
            }}
          >
            {/* Date and Category */}
            <Typography level="body2" sx={{ color: 'text.secondary', mb: '0.5rem' }}>
              <Date dateString={date} /> ⋅ {category}
            </Typography>
            
            {/* Title with default link color */}
            <Typography level="h4" component="h3" sx={{ mb: '0.5rem', fontWeight: 'medium' }}>
              <Link legacyBehavior href={`/blog/${id}`}>
                <a>{title}</a>
              </Link>
            </Typography>
            
            {/* Summary */}
            <Typography level="body2">
              {summary}
            </Typography>
          </ListItem>
        ))}
      </List>
    </Layout>
  );
}

export const getStaticProps: GetStaticProps = async () => {
  const allPostsData = getSortedPostsData();
  return {
    props: {
      allPostsData,
    },
  };
};
