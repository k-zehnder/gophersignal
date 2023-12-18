import Layout from "../components/Layout";
import Typography from "@mui/joy/Typography";
import Avatar from "@mui/joy/Avatar";
import { siteMetaData } from "../lib/siteMetaData";

export default function About() {
  return (
    <Layout>
      <Typography level="h2" sx={{ mb: "1rem" }}>
        About me
      </Typography>
      <Avatar
        sx={{
          "--Avatar-size": "100px",
          mb: "1rem",
        }}
        alt="this person does not exist"
        src={siteMetaData.image}
      >
        JN
      </Avatar>
      <Typography>
        I'm baby waistcoat ugh before they sold out pok pok mlkshk, iceland
        chicharrones. Art party craft beer semiotics Brooklyn, bitters aesthetic
        cornhole authentic vape YOLO food truck waistcoat cliche. Same pork
        belly cornhole wayfarers, hexagon DSA raclette praxis farm-to-table
        edison bulb woke you probably haven't heard of them.
      </Typography>
    </Layout>
  );
}
