# Minimalist Joy UI Blog

## Features

✓ Built with TypeScript

✓ Designed with Joy UI's default styles

✓ Ready to publish with Next.js Markdown blog

✓ Light and dark modes with toggle button

✓ Includes Prettier for code formatting

Created with 💙 by MUI.

## Getting started

[Create a new repository with this template](https://github.com/samuelsycamore/joy-next-blog/generate).

Clone the repo, then run:

```
npm install
```

To start the app in dev mode, run:

```
npm run dev
```

## Customizing

### Personalizing

- Your blog's metadata lives in `/lib/siteMetaData.ts`.
- Your personal contact info is in `/pages/contact.tsx`.

Add your details in these two files.

- Your blog's favicon is located in `/public/`.
- Your avatar (for the About page) is in `/public/images/`.
- The `/public/images/` directory also contains a generic Open Graph card.

Replace these three images with your own.

### Publishing

Blog posts are written in Markdown (`.md`) and kept in the `/posts/` folder.
The demo blog posts contain common Frontmatter keys (`title`, `date`, `summary`, etc.), but you can add or remove as many as you like—just be sure to update `/pages/blog.tsx` and `/pages/blog/[id].tsx` to reflect any changes.

### Styling

This blog uses [Joy UI](https://mui.com/joy-ui/getting-started/overview/), a React component library maintained by [MUI](https://mui.com) (the creators of [Material UI](https://mui.com/material-ui/getting-started/overview/)).

The template is designed almost entirely with Joy UI's default settings, so the customization is minimal out of the box.

There are three ways you can apply custom styles to this template:

- inline, directly on a Joy UI component, using the `sx` prop
- at the theme level, in `/lib/theme.ts`
- globally, on the `<GlobalStyles />` component in `/pages/_app.tsx`

Learn more about [customization approaches in Joy UI](https://mui.com/joy-ui/customization/approaches/).
