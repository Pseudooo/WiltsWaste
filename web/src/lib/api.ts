
const baseUrl = import.meta.env.VITE_API_BASE_URL;

export interface PostcodeAddress {
    address: string;
    uprn: string;
}

export async function fetchAddressesForPostcode(postcode: string): Promise<PostcodeAddress[]> {
    const url = new URL("/addresslist", baseUrl);
    url.searchParams.append("postcode", postcode);

    const response = await fetch(url);
    if (!response.ok) {
        throw new Error("Something go bad");
    }

    const data: PostcodeAddress[] = await response.json();
    return data;
}

export function getWebCalUrl(postcode: string, uprn: string): string {
    const url = new URL("/calendar", baseUrl)
    url.protocol = "webcal";
    url.searchParams.append("postcode", postcode);
    url.searchParams.append("uprn", uprn);

    return url.toString().replace("https://", "webcal://");
}
