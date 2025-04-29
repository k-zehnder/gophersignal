import React from 'react';
import Typography from '@mui/material/Typography';
import FavoriteIcon from '@mui/icons-material/Favorite';

const Footer: React.FC = () => {
  return (
    <footer className="py-4 bg-gray-100 border-t border-gray-300">
      <div className="w-full flex justify-center items-center">
        <Typography
          variant="caption"
          color="textSecondary"
          sx={{
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            gap: '6px',
          }}
        >
          Made with
          <FavoriteIcon
            fontSize="small"
            sx={{
              color: '#ef4444',
              position: 'relative',
              top: '1px',
            }}
          />
          by Gopher Signal
        </Typography>
      </div>
    </footer>
  );
};

export default Footer;
