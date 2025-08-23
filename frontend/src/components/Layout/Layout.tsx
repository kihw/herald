import { Outlet } from 'react-router-dom';
import { Box } from '@mui/material';

// Import layout components when they're created
// import Header from './Header';
// import Sidebar from './Sidebar';
// import Footer from './Footer';

const Layout = () => {
  return (
    <Box sx={{ display: 'flex', flexDirection: 'column', minHeight: '100vh' }}>
      {/* Header will go here */}
      {/* <Header /> */}
      
      <Box sx={{ display: 'flex', flex: 1 }}>
        {/* Sidebar will go here */}
        {/* <Sidebar /> */}
        
        {/* Main content */}
        <Box 
          component="main" 
          sx={{ 
            flex: 1, 
            p: 3,
            bgcolor: 'background.default',
            minHeight: 'calc(100vh - 64px)', // Adjust based on header height
          }}
        >
          <Outlet />
        </Box>
      </Box>
      
      {/* Footer will go here */}
      {/* <Footer /> */}
    </Box>
  );
};

export default Layout;