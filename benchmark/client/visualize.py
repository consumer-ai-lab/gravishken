import pandas as pd
import matplotlib.pyplot as plt

# Read the CSV file
df = pd.read_csv('benchmark_results.csv')

# Create the figure and axis objects
fig, ax = plt.subplots(figsize=(12, 6))

# Plot WebSocket and Polling total times
ax.plot(df['Clients'], df['WebSocket Total (ms)'], marker='o', label='WebSocket Total Time', color='blue')
ax.plot(df['Clients'], df['Polling Total (ms)'], marker='s', label='Polling Total Time', color='red')

# Set labels and title
ax.set_xlabel('Number of Clients')
ax.set_ylabel('Total Time (ms)')
ax.set_title('WebSocket vs Polling: Total Time vs Number of Clients')

# Add legend
ax.legend()

# Add grid for better readability
ax.grid(True, linestyle='--', alpha=0.7)

# Adjust layout and save
plt.tight_layout()
plt.savefig('benchmark_total_time.png')
plt.close()

print("Graph has been saved as 'benchmark_total_time.png'")

# Optional: Print summary statistics
print("\nSummary Statistics:")
print(df[['Clients', 'WebSocket Total (ms)', 'Polling Total (ms)']].describe())

# Optional: Calculate and print the average ratio of Polling to WebSocket total time
avg_ratio = (df['Polling Total (ms)'] / df['WebSocket Total (ms)']).mean()
print(f"\nAverage ratio of Polling to WebSocket total time: {avg_ratio:.2f}")
