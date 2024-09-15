import { Input } from './ui/input';
import { Button } from './ui/button';

export default function AddBatch(){
 

  return (
    <div>
      <h1>Create New Batch</h1>

      <form className="space-y-4">
        <div>
          <label htmlFor="batchName" className="block text-sm font-medium text-gray-700">Batch Name</label>
          <Input type="text" id="batchName" placeholder="Enter batch name" />
        </div>
        <div>
          <label htmlFor="tests" className="block text-sm font-medium text-gray-700">Select Tests</label>
          <select id="tests" multiple className="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm rounded-md">
            <option>Typing Test 1</option>
            <option>Docx Test 1</option>
            <option>Excel Test 1</option>
            <option>Word Test 1</option>
          </select>
        </div>
        <Button type="submit" className="w-full">Create Batch</Button>
      </form>
    </div>
  )
};

