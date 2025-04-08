import { saveAs } from 'file-saver';

const downloadFile = (blob: Blob, filename: string) => {
	saveAs(blob, filename);
};

export default downloadFile;
