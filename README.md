# BubleTimer - Google Cloud Firestore and Cloud Functions

This is an example project from chapter 4 of *Seven Serverless Projects in Seven Weeks* by Sean Johnson.

## Overview

## Usage

## Local Setup

## Technical Design

```json
{
	"user": "4fb61541-4219-41cb-a3c3-3cd525f4d7ab",
	"name": "Albert",
  "email": "albert.camus@combat.org",
	"service_url": "URL",
  "time_categories": [
    {
      "id": "25b69838-1899-11eb-93a1-003ee1cbbd65",
      "name": "Sleeping",
      "color": "08b4ff",
      "active": true
    },
    {
      "id": "b8a7a6f8-ce15-42f6-aa05-988e346f7afb",
      "name": "Writing",
      "color": "ff7bee",
			"active": true
    },
    {
      "id": "92b0f46f-57ca-4831-9629-2315ccc8885d",
      "name": "Thinking",
      "color": "ffc885",
      "active": true
    },
    {
      "id": "135c5eba-a174-46b3-ba0e-8bbcf0035897",
      "name": "Day Job",
      "color": "fffbaa",
      "active": true
    },
    {
      "id": "f54a3fcc-5bcb-44f5-afd9-87b9666c99f9",
      "name": "Reading",
      "color": "abfff7",
      "active": true
    }
  ]
}
```

Firestore schema:

```
users
	user-id
		Document-per-day
			time-slice-index: time-category-uuid
```

Firestore example data:

```
users
	sean@snootymonkey.com
		2020-10-27
			23: "25b69838-1899-11eb-93a1-003ee1cbbd65"
			24: "b8a7a6f8-ce15-42f6-aa05-988e346f7afb"
```

## Participation

Please note that this project is released with a [Contributor Code of Conduct](https://github.com/seven-serverless-projects/bt/blob/mainline/CODE-OF-CONDUCT.md). By participating in this project you agree to abide by its terms.

## License

This code is free software: you can redistribute it and/or modify it under the terms of [The Unlicense](https://unlicense.org/).

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS BE LIABLE FOR ANY CLAIM, DAMAGES OR
OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
OTHER DEALINGS IN THE SOFTWARE.