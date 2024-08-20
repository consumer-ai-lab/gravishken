"use client";

import React, { useState } from "react";

const AdminDashboard = () => {
  const [adminState, setAdminState] = useState({ userName: "Admin" });
  const [file, setFile] = useState(null);
  const [questionInformation, setQuestionInformation] = useState({
    fileType: "",
    timeSlot: "",
    password: "",
    selectedFile: null,
  });

  const handleSubmit = async () => {
    if (file) {
      const data = new FormData();
      data.append("file", file);
      data.append("fileType", questionInformation.fileType);
      data.append("timeSlot", questionInformation.timeSlot);
      data.append("testPassword", questionInformation.password);
    }
  };

  return (
    <div style={styles.container}>
      <h2 className="mb-10">Welcome, {adminState.userName}</h2>

      <div style={styles.formContainer}>
        <h3>Admin Question Uploading Form</h3>

        <div style={styles.formGroup}>
          <label style={styles.label}>Select the Timeslot:</label>
          <select
            style={styles.select}
            onChange={(e) =>
              setQuestionInformation({
                ...questionInformation,
                timeSlot: e.target.value,
              })
            }
          >
            <option value="">Select a Batch</option>
            <option value="Slot 1">Slot 1</option>
            <option value="Slot 2">Slot 2</option>
            <option value="Slot 3">Slot 3</option>
            <option value="Slot 4">Slot 4</option>
          </select>
        </div>

        <div style={styles.formGroup}>
          <label style={styles.label}>Enter the Test Password:</label>
          <input
            type="password"
            style={styles.input}
            onChange={(e) =>
              setQuestionInformation({
                ...questionInformation,
                password: e.target.value,
              })
            }
          />
        </div>

        <div style={styles.formGroup}>
          <label style={styles.label}>Select the Question File Type:</label>
          <select
            style={styles.select}
            onChange={(e) =>
              setQuestionInformation({
                ...questionInformation,
                fileType: e.target.value,
              })
            }
          >
            <option value="">Select the file to upload</option>
            <option value="word">Microsoft Word</option>
            <option value="excel">Microsoft Excel</option>
            <option value="ppt">Microsoft PowerPoint</option>
          </select>
        </div>

        <div style={styles.formGroup}>
          <label style={styles.label}>Upload File:</label>
          <input
            type="file"
          />
        </div>

        <button style={styles.submitButton} onClick={handleSubmit}>
          Upload Data
        </button>
      </div>
    </div>
  );
};

const styles = {
  container: {
    maxWidth: "600px",
    margin: "auto",
    padding: "20px",
    fontFamily: "'Segoe UI', Tahoma, Geneva, Verdana, sans-serif",
  },
  title: {
    textAlign: "center",
    color: "#333",
  },
  logoutButton: {
    display: "block",
    margin: "10px auto",
    padding: "10px 20px",
    backgroundColor: "#f44336",
    color: "#fff",
    border: "none",
    borderRadius: "5px",
    cursor: "pointer",
  },
  formContainer: {
    backgroundColor: "#f9f9f9",
    padding: "20px",
    borderRadius: "10px",
    boxShadow: "0 0 10px rgba(0,0,0,0.1)",
  },
  formTitle: {
    marginBottom: "20px",
    textAlign: "center",
    color: "#444",
  },
  formGroup: {
    marginBottom: "15px",
  },
  label: {
    display: "block",
    marginBottom: "5px",
    color: "#555",
  },
  input: {
    width: "100%",
    padding: "10px",
    border: "1px solid #ccc",
    borderRadius: "5px",
  },
  select: {
    width: "100%",
    padding: "10px",
    border: "1px solid #ccc",
    borderRadius: "5px",
  },
  fileInput: {
    width: "100%",
    padding: "10px",
  },
  submitButton: {
    display: "block",
    width: "100%",
    padding: "10px",
    backgroundColor: "#4CAF50",
    color: "#fff",
    border: "none",
    borderRadius: "5px",
    cursor: "pointer",
  },
};

export default AdminDashboard;
