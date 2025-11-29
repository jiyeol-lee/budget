import { describe, it, expect } from 'vitest';
import { extractErrorMessage } from './api';

/**
 * Test suite for extractErrorMessage utility function
 *
 * Tests cover:
 * - Valid JSON with string error field (standard backend responses)
 * - Valid JSON without error field or with non-string error field
 * - Invalid JSON (plain text, HTML)
 * - Empty or falsy input handling
 */

describe('extractErrorMessage', () => {
	describe('with valid JSON containing error field', () => {
		it('extracts error message from standard backend response', () => {
			const input = '{"error": "Budget for this month/year already exists"}';
			expect(extractErrorMessage(input)).toBe('Budget for this month/year already exists');
		});

		it('extracts error message with special characters', () => {
			const input = '{"error": "Value must be > 0 and < 100"}';
			expect(extractErrorMessage(input)).toBe('Value must be > 0 and < 100');
		});

		it('extracts empty error message', () => {
			const input = '{"error": ""}';
			expect(extractErrorMessage(input)).toBe('');
		});

		it('extracts error message with unicode characters', () => {
			const input = '{"error": "예산이 이미 존재합니다"}';
			expect(extractErrorMessage(input)).toBe('예산이 이미 존재합니다');
		});

		it('extracts error message with newlines and whitespace', () => {
			const input = '{"error": "Line 1\\nLine 2"}';
			expect(extractErrorMessage(input)).toBe('Line 1\nLine 2');
		});

		it('extracts error message when other fields are present', () => {
			const input = '{"error": "Not found", "code": 404, "timestamp": "2024-01-01"}';
			expect(extractErrorMessage(input)).toBe('Not found');
		});
	});

	describe('with valid JSON without error field', () => {
		it('returns original text when error field is missing', () => {
			const input = '{"message": "Something went wrong"}';
			expect(extractErrorMessage(input)).toBe(input);
		});

		it('returns original text for empty JSON object', () => {
			const input = '{}';
			expect(extractErrorMessage(input)).toBe(input);
		});

		it('returns original text when error field is not a string', () => {
			const input = '{"error": 123}';
			expect(extractErrorMessage(input)).toBe(input);
		});

		it('returns original text when error field is null', () => {
			const input = '{"error": null}';
			expect(extractErrorMessage(input)).toBe(input);
		});

		it('returns original text when error field is an object', () => {
			const input = '{"error": {"code": "INVALID"}}';
			expect(extractErrorMessage(input)).toBe(input);
		});

		it('returns original text when error field is an array', () => {
			const input = '{"error": ["error1", "error2"]}';
			expect(extractErrorMessage(input)).toBe(input);
		});

		it('returns original text when error field is a boolean', () => {
			const input = '{"error": true}';
			expect(extractErrorMessage(input)).toBe(input);
		});

		it('returns original text for JSON array', () => {
			const input = '[{"id": 1}, {"id": 2}]';
			expect(extractErrorMessage(input)).toBe(input);
		});
	});

	describe('with invalid JSON (plain text)', () => {
		it('returns plain text error message as-is', () => {
			const input = 'Internal Server Error';
			expect(extractErrorMessage(input)).toBe('Internal Server Error');
		});

		it('returns HTML error response as-is', () => {
			const input = '<html><body>Error</body></html>';
			expect(extractErrorMessage(input)).toBe(input);
		});

		it('returns partial JSON as-is', () => {
			const input = '{"error": "incomplete';
			expect(extractErrorMessage(input)).toBe(input);
		});

		it('returns random text as-is', () => {
			const input = 'Something went wrong: connection refused';
			expect(extractErrorMessage(input)).toBe(input);
		});

		it('returns whitespace-only text as-is', () => {
			const input = '   ';
			expect(extractErrorMessage(input)).toBe(input);
		});

		it('returns numeric string as-is', () => {
			const input = '500';
			expect(extractErrorMessage(input)).toBe(input);
		});
	});

	describe('with empty or falsy input', () => {
		it('returns "Unknown error" for empty string', () => {
			expect(extractErrorMessage('')).toBe('Unknown error');
		});
	});
});
