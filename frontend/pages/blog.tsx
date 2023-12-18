import Layout from "../components/Layout";
import { getSortedPostsData } from "../lib/posts";
import Link from "next/link";
import Date from "../components/Date";
import { GetStaticProps } from "next";
import Typography from "@mui/joy/Typography";
import List from "@mui/joy/List";
import ListItem from "@mui/joy/ListItem";

export default function Blog({
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
      <Typography level="h2">Posts</Typography>
      <List
        sx={{
          display: "flex",
          gap: "10px",
        }}
      >
        {allPostsData.map(({ id, date, category, title, summary }) => (
          <ListItem
            key={id}
            sx={{
              display: "flex",
              flexDirection: "column",
              alignItems: "flex-start",
            }}
          >
            <Typography level="body3">
              <Date dateString={date} /> ⋅ {category}
            </Typography>
            <Typography level="h4" component="p">
              <Link href={`/blog/${id}`}>{title}</Link>
            </Typography>
            {summary}
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
