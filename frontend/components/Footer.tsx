import React from 'react';
import FavoriteIcon from '@mui/icons-material/Favorite';
import Typography from '@mui/material/Typography';

const Footer: React.FC = () => {
  return (
    <footer className="py-4 bg-gray-100 border-t border-gray-300">
      <div
        className="w-full flex justify-center items-center"
        style={{ maxWidth: '1200px', margin: '0 auto' }}
      >
        <Typography
          variant="body2"
          className="text-gray-600"
          component="p"
          sx={{
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            textAlign: 'center',
          }}
        >
          Made with{' '}
          <FavoriteIcon fontSize="small" sx={{ color: 'red', mx: 0.5 }} /> by{' '}
          <a
            href="https://github.com/k-zehnder/gophersignal"
            target="_blank"
            rel="noopener noreferrer"
            style={{
              color: '#007bff',
              textDecoration: 'none',
              fontWeight: 'bold',
              marginLeft: '0.25rem',
            }}
            onMouseEnter={(e) =>
              (e.currentTarget.style.textDecoration = 'underline')
            }
            onMouseLeave={(e) =>
              (e.currentTarget.style.textDecoration = 'none')
            }
          >
            GopherSignal
          </a>
        </Typography>
      </div>
    </footer>
  );
};

export default Footer;
