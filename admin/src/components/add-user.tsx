import React, { useState } from 'react';
import { User } from '@common/types';
import { Input } from './ui/input';
import { Button } from './ui/button';

export function AddUser() {
  

  return (
    <div>
      <h1>Add Users from CSV</h1>

      <form className="space-y-4">
        <div>
          <label htmlFor="csvFile" className="block text-sm font-medium text-gray-700">Upload CSV File</label>
          <Input type="file" id="csvFile" accept=".csv" />
        </div>
        <Button type="submit" className="w-full">Upload and Add Users</Button>
      </form>

    </div>
  )
};
