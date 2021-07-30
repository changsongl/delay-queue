package encode

import (
	"github.com/changsongl/delay-queue/job"
	"github.com/changsongl/delay-queue/pkg/lock"
	"github.com/stretchr/testify/require"
	"testing"
)

var longText = `
The Bible is not a single book but a collection of books, whose complex development is not completely understood. The books began as songs and stories orally transmitted from generation to generation before being written down in a process that began sometime around the start of the first millennium BCE and continued for over a thousand years. The Bible was written and compiled by many people, from a variety of disparate cultures, most of whom are unknown.[18] British biblical scholar John K. Riches wrote:[19]
[T]he biblical texts were produced over a period in which the living conditions of the writers – political, cultural, economic, and ecological – varied enormously. There are texts which reflect a nomadic existence, texts from people with an established monarchy and Temple cult, texts from exile, texts born out of fierce oppression by foreign rulers, courtly texts, texts from wandering charismatic preachers, texts from those who give themselves the airs of sophisticated Hellenistic writers. It is a time-span which encompasses the compositions of Homer, Plato, Aristotle, Thucydides, Sophocles, Caesar, Cicero, and Catullus. It is a period which sees the rise and fall of the Assyrian empire (twelfth to seventh century) and of the Persian empire (sixth to fourth century), Alexander's campaigns (336–326), the rise of Rome and its domination of the Mediterranean (fourth century to the founding of the Principate, 27 BCE), the destruction of the Jerusalem Temple (70 CE), and the extension of Roman rule to parts of Scotland (84 CE).
Hebrew Bible from 1300. page 20, Genesis.
Hebrew Bible from 1300. Genesis.
Considered to be scriptures (sacred, authoritative religious texts), the books were compiled by different religious communities into various biblical canons (official collections of scriptures). The earliest compilation, containing the first five books of the Bible and called the Torah (meaning "law", "instruction", or "teaching") or Pentateuch ("five books"), was accepted as Jewish canon by the 5th century BCE. A second collection of narrative histories and prophesies, called the Nevi'im ("prophets"), was canonized in the 3rd century BCE. A third collection called the Ketuvim ("writings"), containing psalms, proverbs, and narrative histories, was canonized sometime between the 2nd century BCE and the 2nd century CE. These three collections were written mostly in Hebrew, with some parts in Aramaic, and together form the Hebrew Bible or "TaNaKh" (a portmanteau of "Torah", "Nevi'im", and "Ketuvim").[20]
Greek-speaking Jews in Alexandria and elsewhere in the Jewish diaspora considered additional scriptures, composed between 200 BCE and 100 CE and not included in the Hebrew Bible, to be canon. These additional texts were included in a translation of the Hebrew Bible into Koine Greek (common Greek spoken by ordinary people) known as the Septuagint (meaning "the work of the seventy"), which began as a translation of the Torah made around 250 BCE and continued to develop for several centuries. The Septuagint contained all of the books of the Hebrew Bible, re-organized and with some textual differences, with the additional scriptures interspersed throughout.[21]
Saint Paul Writing His Epistles, 16th-century painting.
During the rise of Christianity in the 1st century CE, new scriptures were written in Greek about the life and teachings of Jesus Christ, who Christians believed was the messiah prophesized in the books of the Hebrew Bible. Two collections of these new scriptures – the Pauline epistles and the Gospels – were accepted as canon by the end of the 2nd century CE. A third collection, the catholic epistles, were canonized over the next few centuries. Christians called these new scriptures the "New Testament", and began referring to the Septuagint as the "Old Testament".[22]
Between 385 and 405 CE, the early Christian church translated its canon into Vulgar Latin (the common Latin spoken by ordinary people), a translation known as the Vulgate, which included in its Old Testament the books that were in the Septuagint but not in the Hebrew Bible. The Vulgate introduced stability to the Bible, but also began the East-West Schism between Latin-speaking Western Christianity (led by the Catholic Church) and multi-lingual Eastern Christianity (led by the Eastern Orthodox Church). Christian denominations' biblical canons varied not only in the language of the books, but also in their selection, organization, and text.[23]
Jewish rabbis began developing a standard Hebrew Bible in the 1st century CE, maintained since the middle of the first millennium by the Masoretes, and called the Masoretic Text. Christians have held ecumenical councils to standardize their biblical canon since the 4th century CE. The Council of Trent (1545–63), held by the Catholic Church in response to the Protestant Reformation, authorized the Vulgate as its official Latin translation of the Bible. The Church deemed the additional books in its Old Testament that were interspersed among the Hebrew Bible books to be "deuterocanonical" (meaning part of a second or later canon). Protestant Bibles either separated these books into a separate section called the "Apocrypha" (meaning "hidden away") between the Old and New Testaments, or omitted them altogether. The 17th-century Protestant King James Version was the most ubiquitous English Bible of all time, but it has largely been superseded by modern translations.[24]`

func runEncodeTest(t *testing.T, encoder Encoder) {
	for _, j := range testCases(t) {
		str, err := encoder.Encode(j)
		require.NoError(t, err)

		jDecode := &job.Job{}
		err = encoder.Decode(str, jDecode)
		t.Logf("%+v", jDecode)

		require.NoError(t, err)
		require.Equal(t, j.ID, jDecode.ID)
		require.Equal(t, j.TTR, jDecode.TTR)
		require.Equal(t, j.Delay, jDecode.Delay)
		require.Equal(t, j.Topic, jDecode.Topic)
		require.True(t, j.Version.Equal(jDecode.Version))
		require.Equal(t, j.Body, jDecode.Body)
		require.Equal(t, j.Version.String(), jDecode.Version.String())
	}
}

func testCases(t *testing.T) []*job.Job {
	jEmpty, err := job.New("jobTopic", "1sdsa", 0, 0, "", func(name string) lock.Locker {
		return nil
	})
	require.NoError(t, err)

	jAllEmpty := &job.Job{}

	jShort, err := job.New("jobTopicdsad", "1sdsadsads", 10,
		20, "gfdgsdas", func(name string) lock.Locker {
			return nil
		},
	)
	require.NoError(t, err)

	jLong, err := job.New(job.Topic(longText), job.ID(longText), 10000000000,
		10000000000, job.Body(longText), func(name string) lock.Locker {
			return nil
		},
	)
	require.NoError(t, err)

	return []*job.Job{
		jEmpty,
		jAllEmpty,
		jShort,
		jLong,
	}
}
