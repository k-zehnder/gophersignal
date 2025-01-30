import React from 'react';
import Typography from '@mui/material/Typography';
import { FaSkullCrossbones } from 'react-icons/fa';
import { useTheme } from '@mui/material/styles';

const Footer: React.FC = () => {
  const theme = useTheme();

  return (
    <footer className="py-4 bg-gray-100 border-t border-gray-300">
      <div className="w-full flex flex-col justify-center items-center">
        <Typography
          variant="body2"
          className="text-gray-600"
          sx={{
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            fontSize: '0.95rem',
            fontWeight: 500,
            letterSpacing: '0.5px',
          }}
        >
          <span
            style={{
              display: 'flex',
              alignItems: 'center',
              gap: '8px',
            }}
          >
            Made Possible
            <FaSkullCrossbones
              style={{
                fontSize: '1.6rem',
                color: theme.palette.mode === 'dark' ? '#fff' : '#000',
                position: 'relative',
                top: '1px',
              }}
            />
            by Uncle Dennis
          </span>
        </Typography>
      </div>
    </footer>
  );
};

export default Footer;
