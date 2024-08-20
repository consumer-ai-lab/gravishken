import React, { useState, useEffect } from 'react';// Adjust the import based on your project structure
import AdminNavbar from './Navbar';

const BASE_URL = 'your_base_url_here'; // Replace with actual base URL

interface BatchData {
  [0]: string;
  [1]: string | null;
  [2]: string;
}

interface AdminState {
  userName: string;
}

const AdminViewSubmission: React.FC = () => {
  const [batchData, setBatchData] = useState<BatchData[]>([]);
  const [batch, setBatch] = useState<string>("");
  const [adminState, setAdminState] = useState<AdminState>({ userName: '' });
  const [loginState, setLoginState] = useState<'checking' | 'loggedIn' | 'notLoggedIn' | 'networkError'>('checking');
  const [openSuccessAlert, setOpenSuccessAlert] = useState<boolean>(false);
  const [progressAlert, setProgressAlert] = useState<boolean>(false);

  
  const StaticTable: React.FC = () => (
    <div className="container" style={{ width: '70%', margin: 'auto' }}>
      <p style={{ textAlign: 'center', fontWeight: 'bolder' }}>
        Select the Preferred Slot and click the Fetch Batch Data Button
      </p>
    </div>
  );

  const DynamicTable: React.FC<{ batchData: BatchData[] }> = ({ batchData }) => (
    <div className="container" style={{ width: '70%', margin: 'auto' }}>
      <table>
        <thead>
          <tr>
            <th style={{ fontWeight: 'bold', fontSize: '18px' }}>
              <center>Roll Number</center>
            </th>
            <th style={{ fontWeight: 'bold', fontSize: '18px' }}>
              <center>Download Single File</center>
            </th>
            <th style={{ fontWeight: 'bold', fontSize: '18px' }}>
              <center>Exam Attendance</center>
            </th>
          </tr>
        </thead>
        <tbody>
          {batchData.map((b, index) => (
            <tr key={index}>
              <td>
                <center>{b[0]}</center>
              </td>
              <td>
                {b[2] === "Present" ? (
                  <center>
                    <button>
                      Download
                    </button>
                  </center>
                ) : (
                  <center style={{ color: 'red', fontWeight: 'bold' }}>-</center>
                )}
              </td>
              <td>
                {b[2] === "Present" ? (
                  <center style={{ color: 'green', fontWeight: 'bold' }}>{b[2]}</center>
                ) : (
                  <center style={{ color: 'red', fontWeight: 'bold' }}>{b[2]}</center>
                )}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );

  return (
    <div>
      <AdminNavbar />
      <div className="container">
        
        </div>
    </div>
  );
};

export default AdminViewSubmission;
