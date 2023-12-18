import Layout from "../../components/Layout";
import { getAllPostIds, getPostData } from "../../lib/posts";
import Head from "next/head";
import Date from "../../components/Date";
import { GetStaticProps, GetStaticPaths } from "next";
import Typography from "@mui/joy/Typography";

export default function Post({
  postData,
}: {
  postData: {
    title: string;
    summary: string;
    category: string;
    author: string;
    date: string;
    contentHtml: string;
  };
}) {
  return (
    <Layout>
      <Head>
        <title>{postData.title}</title>
      </Head>

      <Typography level="display2" fontSize="30px" mb={1}>
        {postData.title}
      </Typography>
      <Typography level="body2">
        {postData.author} ⋅ <Date dateString={postData.date} /> ⋅{" "}
        {postData.category}
      </Typography>
      <div dangerouslySetInnerHTML={{ __html: postData.contentHtml }} />
    </Layout>
  );
}

export const getStaticPaths: GetStaticPaths = async () => {
  const paths = getAllPostIds();
  return {
    paths,
    fallback: false,
  };
};

export const getStaticProps: GetStaticProps = async ({ params }) => {
  const postData = await getPostData(params?.id as string);
  return {
    props: {
      postData,
    },
  };
};
